'use strict';

angular.module('myApp.playList', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
    $routeProvider.when('/play-list', {
        templateUrl: 'web/app/components/play-list/play-list.html',
        controller: 'PlayListCtrl'
    });
}])

.controller('PlayListCtrl', ["$scope", "$http", "$location", "Data", function($scope, $http, $location, Data) {
    $scope.tracks = [];

    var httpRequest = $http({
        method: 'GET',
        url: '/api/play-list/'
    }).success(function(data, status) {
        if (data.Code != 500) {
            $scope.tracks = data;
        }
    });

    $scope.clearPlaylist = function() {
        var httpRequest = $http({
            method: 'GET',
            url: '/api/clear/'
        }).success(function(data, status) {
            $location.path('/library');
        });

    };

    $scope.$watch(function () { return Data.getState(); }, function (state) {
        if (state) {
            $("td").css("color", "#333");
            $("td:contains('" + state.File + "')").css("color", "#d43943");
        }
    });

}]);