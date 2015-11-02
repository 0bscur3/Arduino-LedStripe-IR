var app = angular.module('LedServer', ['ngMaterial', 'ngWebsocket']);

app.config(function($mdThemingProvider) {
    $mdThemingProvider.theme('default')
        .primaryPalette('blue')
        .accentPalette('blue');
});

app.controller('AppCtrl', ['$scope', '$websocket', function($scope, $websocket){

    var ws = $websocket.$new({
        url: "ws://localhost:8081",
        protocols: []
    });

    ws.$on('$open', function () {
        console.log('Websocket connected');

        ws.$emit("color", [255, 255, 255]);
        ws.$emit("cmd","FADE");
    });

    ws.$on('$close', function () {
        console.log('Websocket closed');
    });

    $scope.color = {
        red: 255,
        green: 100,
        blue: 100
    };

    $scope.messages = [
        'Hello',
        'Hello2',
        'Received Message'
    ]
}]);