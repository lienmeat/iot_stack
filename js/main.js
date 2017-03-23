$(document).ready(function() {
	getDeviceID();
});

var ws_mqtt_host = "ericslien.com:9883";

var device_type = "bottlelamp";

var device_id = "";

var topics = {
	"publish": {
		"topic_set_mode": device_type + "/{device_id}/set/mode",
		"topic_set_speed": device_type + "/{device_id}/set/speed",
		"topic_set_color": device_type + "/{device_id}/set/color",
		"topic_statusupdate": device_type + "/{device_id}/statusupdate"
	},
	"subscribe": {
		"topic_mode": device_type + "/{device_id}/mode",
		"topic_speed": device_type + "/{device_id}/speed",
		"topic_color": device_type + "/{device_id}/color",
		"topic_modelist": device_type + "/{device_id}/modelist"
	}
}

var current_mode = "loading";

var mqttc = mqtt.connect('ws://' + ws_mqtt_host);

mqttc.on('connect', function () {
	connectMQTT();
})

mqttc.on('message', function (topic, message) {
	message = message.toString()
	console.log(message)
	switch (topic) {
		case topics.subscribe.topic_mode:
			onMode(message);
			break;
		case topics.subscribe.topic_speed:
			onSpeed(message);
			break;
		case topics.subscribe.topic_color:
			onColor(message);
			break;
		case topics.subscribe.topic_modelist:
			onModelist(message);
			break;
	}
})

mqttc.on('error', function(err) {
	console.log(err)
})

$("#set_device_name").click(setupDevice);

$("#conf_panel_open").click(function() {
	$("#device_connect").show();
	$("#lampcontrol").hide();
	$("#conf_panel_open").hide();
});

$("#device_name").keypress(function(event) {
	if(event.keycode == 13) {
		setupDevice(e);
	}
});

function getDeviceID() {
	if( device_id.length <= 0 ) {
		var id = localStorage.getItem("device_id");
		if( id && id.length > 0 ) {
			device_id = id;
			$("#device_name").val(device_id);
			connectMQTT();
			return id;
		}
		return device_id;
	}
}

function setDeviceID(id) {
	if(id && id.length > 0) {
		device_id = id;
		localStorage.setItem("device_id", id);
	}
}

function setupDevice(e) {
	var dev_name = $("#device_name").val();
	console.log(dev_name);
	if(dev_name.length > 0) {
		console.log("setting device id");
		setDeviceID(dev_name);
	}
	if(device_id.length > 0) {
		for(var i in topics.publish) {
			topics.publish[i] = topics.publish[i].replace("{device_id}", device_id);
		}
		for(var i in topics.subscribe) {
			topics.subscribe[i] = topics.subscribe[i].replace("{device_id}", device_id);
			mqttc.subscribe(topics.subscribe[i]);
		}
		getStatusUpdate();
		$("#device_connect").hide();
		$("#lampcontrol").show();
		$("#conf_panel_open").show();
	}
}

function connectMQTT() {
	mqttc.publish('presence', 'Hello mqtt');
	setupDevice();
}

function onMode(message) {
	$('#mode').val(message);
	renderColorCtrl();
}

function onSpeed(message) {
	$('#speed').val(parseInt(message));
}

function onColor(message) {
	$('#color').val(message);
}

function onModelist(message) {
	var options = message.split(',');
	renderModeOptions(options, current_mode);
}

function sendMode() {
	renderColorCtrl();
	var mode = $('#mode').val()
	console.log("mode:" + mode);
	mqttc.publish(topics.publish.topic_set_mode, mode);
}

function sendSpeed() {
	var speed = parseInt($('#speed').val());
	$('#speeddisp').html(speed);
	console.log("speed:" + speed);
	mqttc.publish(topics.publish.topic_set_speed, speed);
}

function sendColor() {
	var color = $('#color').val();
	console.log("color:" + color);
	mqttc.publish(topics.publish.topic_set_color, color);
}

function getStatusUpdate() {
	if(device_id && device_id.length != 0) {
		mqttc.publish(topics.publish.topic_statusupdate, 1);
	}
}

function renderModeOptions(options, current) {
	var opts = "";
	for(var i in options) {
		var sel = "";
		if(options[i] == current) {
			sel = " selected=\"selected\"";
		}
		opts += "<option value=\"" + options[i] + "\" " + sel + ">" + options[i] + "</option>";
	}
	$('#mode').html(opts);
	renderColorCtrl();
}

function renderColorCtrl() {
	if($('#mode option:selected').text() == "Solid Color") {
		$('#colorselection').removeClass('hidden');
	}
	else{
		$('#colorselection').addClass('hidden');;
	}
}