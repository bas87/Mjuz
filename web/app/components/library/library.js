'use strict';

angular.module('myApp.library', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
    $routeProvider.when('/library', {
        templateUrl: 'web/app/components/library/library.html',
        controller: 'LibraryCtrl'
    });
}])

.controller('LibraryCtrl', ["$scope", "$http", "$base64", function($scope, $http, $base64) {
	$scope.url = null;
    $scope.tracks = [];

    var httpRequest = $http({
        method: 'GET',
        url: '/api/list/'
    }).success(function(data, status) {
        if (data.Code != 500) {
            $scope.tracks = data;
        }
    });

    $scope.addToPlaylist = function(data) {
        var httpRequest = $http({
            method: 'GET',
            url: '/api/append/?path=' + $base64.encode(data.Path)
        }).success(function(data, status) {
            alert("Added to playlist!");
        });
    };

    $scope.download = function () {
        var httpRequest = $http({
            method: 'GET',
            url: '/api/download/?url=' + $base64.encode($scope.url)
        }).success(function(data, status) {
            alert("Complete!");
        });
    };

}]);