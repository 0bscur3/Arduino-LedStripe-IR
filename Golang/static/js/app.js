var app = angular.module('LedServer', ['ngMaterial', 'ngWebSocket']);

app.config(['$mdThemingProvider', function($mdThemingProvider) {
    $mdThemingProvider.theme('default')
        .primaryPalette('blue')
        .accentPalette('blue');
}]);

app.controller('AppCtrl', ['$scope', '$websocket', function($scope, $websocket){

    var ws = $websocket("ws://localhost:8081/");

    ws.onMessage(function(message) {
        var data = JSON.parse(message.data);
        console.log(data);
        if (data.event === "cmd") {
            handleCommand(data.data);
        } else if (data.event === "color") {
            handleColor(data.data);
        }

    });

    function handleCommand (command) {

        var validCommands = ['POWER_ON', 'POWER_OFF', 'FADE', 'FLASH'];
        var value = 0;

        if(validCommands.indexOf(command) < 0){
            return;
        }

        if (command.slice(-2) === "ON") {
            value = true;
        } else if (command.slice(-3) === "OFF") {
            value = false;
        }

        var parts = command.split("_");
        var option = parts[0];

        if (option === "FADE") {
            $scope.options.fade = !$scope.options.fade;
        } else if (option === "FLASH") {
            $scope.options.flash = !$scope.options.flash;
        } else if (option === "POWER") {
            $scope.status.power = value;
        }


    }

    function handleColor (colors) {
        $scope.color.red = colors[0];
        $scope.color.green = colors[1];
        $scope.color.blue = colors[2];
    }

    $scope.status = {
        power: false
    };

    $scope.options = {
        fade: false,
        flash: false
    };

    $scope.color = {
        red: 255,
        green: 100,
        blue: 100
    };

    $scope.messages = [
        'Hello',
        'Hello2',
        'Received Message'
    ];

    $scope.sendColor = function(red, green, blue) {
        ws.send({event:'color', data:[red, green, blue]});
    };

    $scope.sendPower = function(value) {
        if (value === true) {
            $scope.sendCommand('POWER_ON');
        } else {
            $scope.sendCommand('POWER_OFF');
        }
    };

    $scope.sendCommand = function(command) {
/*
        command += "_";
        console.log(value);

        if (value === true) {
            command += "ON";
        } else {
            command += "OFF";
        }*/

        ws.send({
            event:'cmd',
            data:command
        });
    };
}]);