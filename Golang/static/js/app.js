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

        if (data.event === "cmd") {
            handleCommand(data.data);
        } else if (data.event === "color") {
            handleColor(data.data);
        }

    });

    function handleCommand (command) {

        var valid = false;
        var value = 0;

        if (command.slice(-2) === "ON") {
            valid = true;
            value = true;
        } else if (command.slice(-3) === "OFF") {
            valid = true;
            value = false;
        }

        if (!valid) {
            return;
        }

        var parts = command.split("_");
        var option = parts[0];

        if (option === "FADE") {
            $scope.options.fade = value;
        } else if (option === "FLASH") {
            $scope.options.flash = value;
        } else if (option === "STROBE") {
            $scope.options.strobe = value;
        }


    }

    function handleColor (colors) {
        $scope.color.red = colors[0];
        $scope.color.green = colors[1];
        $scope.color.blue = colors[2];
    }

    $scope.options = {
        fade: true,
        flash: true,
        strobe: true
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

    $scope.sendCommand = function(command, value) {

        command += "_";
        console.log(value);

        if (value === true) {
            command += "ON";
        } else {
            command += "OFF";
        }

        ws.send({
            event:'cmd',
            data:command
        });
    };
}]);