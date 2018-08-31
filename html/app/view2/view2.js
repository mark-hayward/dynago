'use strict';

angular.module('myApp.view2', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/view2', {
    templateUrl: 'view2/view2.html',
    controller: 'View2Ctrl'
  });
}])

.controller('View2Ctrl',  ['$scope', function($scope) {
        $scope.regions= {
            euwest1: "Ireland",
            euwest2: "London",
            useast1: "North Virginia",
            uswest2: "Oregan"
        };
        $scope.dynamotable= {
            text: 'DynamoDBTable',
            word: /^\s*\w*\s*$/
        };
}]);