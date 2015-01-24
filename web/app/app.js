'use strict';

angular.module('myApp', [
    'ngRoute',
    'readableTime',
    'base64',
    'myApp.library',
    'myApp.playList'
]).
config(['$routeProvider', function($routeProvider) {
    $routeProvider.otherwise({
        redirectTo: '/library'
    });
}])
.factory('Data', function(){
    var data =
        {
            state: null
        };
    
    return {
        getState: function () {
            return data.state;
        },
        setState: function (state) {
            data.state = state;
        }
    };
})
.controller('PlayerCtrl', ["$scope", "$http", "$timeout", "Data", function($scope, $http, $timeout, Data) {

    $scope.play = function() {
        var httpRequest = $http({
            method: 'GET',
            url: '/api/play/'
        }).success(function(data, status) {
            console.info("Let's play!");
        });

    };


    $scope.stop = function() {
        var httpRequest = $http({
            method: 'GET',
            url: '/api/stop/'
        }).success(function(data, status) {
            console.info("Let's stop!");
        });

    };


    $scope.pause = function() {
        var httpRequest = $http({
            method: 'GET',
            url: '/api/toogle-pause/'
        }).success(function(data, status) {
            console.info("Timeout!");
        });

    };


    $scope.prev = function() {
        var httpRequest = $http({
            method: 'GET',
            url: '/api/prev/'
        }).success(function(data, status) {
            console.info("Let's to prev song!");
        });

    };


    $scope.next = function() {
        var httpRequest = $http({
            method: 'GET',
            url: '/api/next/'
        }).success(function(data, status) {
            console.info("Let's to next song!");
        });

    };

    $scope.$watch(function () { return Data.getState(); }, function (state) {
        $scope.state = state;
    });


    (function gestState() {
        $timeout(function() {
            var httpRequest = $http({
                method: 'GET',
                url: '/api/info/'
            }).success(function(data, status) {
                Data.setState(data);
            });

            gestState();
        }, 1000);
    })();


}]).controller('HeaderCtrl', ["$scope", "$location", function($scope, $location) {

    $scope.isActive = function (viewLocation) { 
        return viewLocation === $location.path();
    };

}]);