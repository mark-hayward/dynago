'use strict';

angular.module('myApp.view3', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/view3', {
    templateUrl: 'view3/view3.html',
    controller: 'View3Ctrl'
  });
}])

.controller('View3Ctrl', ['$scope', function($scope) {
    $scope.s3bucketname= {
        text: 'S3 Bucket Name',
        word: /^\s*\w*\s*$/
    };
    $scope.s3key= {
        text: 'S3 Key Name',
        word: /^\s*\w*\s*$/
    };

}]);