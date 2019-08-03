// SPDX-License-Identifier: Apache-2.0

'use strict';

var app = angular.module('application', []);

// Angular Controller
app.controller('appController', function($scope, appFactory){

	$("#success_owner").hide();
	$("#success_create").hide();
	$("#error_owner").hide();
	$("#error_query").hide();
	
	$scope.queryAllProduce = function(){

		appFactory.queryAllProduce(function(data){
			var array = [];
			for (var i = 0; i < data.length; i++){
				parseInt(data[i].Key);
				data[i].Record.Key = parseInt(data[i].Key);
				array.push(data[i].Record);
			}
			array.sort(function(a, b) {
			    return parseFloat(a.Key) - parseFloat(b.Key);
			});
			$scope.all_produce = array;
		});
	}

	$scope.queryProduce = function(){

		var id = $scope.produce_id;

		appFactory.queryProduce(id, function(data){
			$scope.query_produce = data;

			if ($scope.query_produce == "Could not locate produce"){
				console.log()
				$("#error_query").show();
			} else{
				$("#error_query").hide();
			}
		});
	}

	$scope.recordProduce = function(){

		appFactory.recordProduce($scope.produce, function(data){
			$scope.create_produce = data;
			$("#success_create").show();
		});
	}

	$scope.changeOwner = function(){

		appFactory.changeOwner($scope.owner , function(data){
			$scope.change_owner = data;
			if ($scope.change_owner == "Error: no produce catch found"){
				$("#error_owner").show();
				$("#success_owner").hide();
			} else{
				$("#success_owner").show();
				$("#error_owner").hide();
			}
		});
	}

});

// Angular Factory
app.factory('appFactory', function($http){
	
	var factory = {};

    factory.queryAllProduce = function(callback){

    	$http.get('/get_all_produce/').success(function(output){
			callback(output)
		});
	}

	factory.queryProduce = function(id, callback){
    	$http.get('/get_produce/'+id).success(function(output){
			callback(output)
		});
	}

	factory.recordProduce = function(data, callback){

		data.location = data.longitude + ", "+ data.latitude;

		var produce = data.id + "-" + data.location + "-" + data.timestamp + "-" + data.owner  + "-" + data.vessel;

    	$http.get('/add_produce/'+produce).success(function(output){
			callback(output)
		});
	}

	factory.changeOwner = function(data, callback){

		var owner  = data.id + "-" + data.name;

    	$http.get('/change_owner/'+owner ).success(function(output){
			callback(output)
		});
	}

	return factory;
});


