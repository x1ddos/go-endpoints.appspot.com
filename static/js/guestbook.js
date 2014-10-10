(function(window) {
  'use strict';
  var angular = window.angular;

  /**
   * AngularJS app
   */
  var app = angular.module("guestbook", []).
    constant('ClientID', window.CLIENT_ID).
    constant('Scopes', window.SCOPES).
    constant('ApiRoot', '//' + window.location.host + '/_ah/api').
    provider('API', ApiProvider);

  /**
   * Main and the only controller of the guestbook app
   */
  function GreetingCtrl($log, $scope, API) {
    $scope.greetings = API.listGreetings({limit: 50});
    $scope.newgreet = {};
    $scope.user = null;

    /**
     * Adds a new message to the greets list in memory and submits it
     * to the backend meanwhile.
     */
    $scope.signGuestbook = function() {
      if ($scope.newgreet.content.match(/^\s*$/)) return;

      // add locally
      angular.extend($scope.newgreet, {
        'date': new Date(),
        'author': $scope.user ? $scope.user.given_name : 'Anonymous User'
      });
      var len = $scope.greetings.unshift($scope.newgreet),
          data = {message: $scope.newgreet.content};

      // store on the backend
      API.signGuestbook(data).then(
        function(resp){
          var greet = resp.data,
              idx = $scope.greetings.length - len;
          // replace greet mock with the real one from the backend
          $scope.greetings[idx] = greet;
        },
        function(){
          // remove the message if we failed to store it on the backend
          var idx = $scope.greetings.length - len,
              greet = $scope.greetings.splice(idx, 1)[0];
          angular.copy(greet, $scope.newgreet);
        });

      // reset the form
      // TODO(alex): form.setPrisine() ?
      $scope.newgreet = {};
    }

    /**
     * Logs a user in standard (non-immediate) mode when clicked on
     * "Login" button.
     */
    $scope.login = function() {
      API.login(false);
    }

    /**
     * Watch the status of user authentication.
     */
    $scope.$on('Authenticated', function($event, authed){
      $scope.user = authed ? API.getUserInfo() : null;
    })
  }

  /**
   * Our API client provider.
   */
  function ApiProvider($httpProvider) {
    var gapi, authToken;

    /**
     * This interceptor adds "Authorization" header when auth info is available.
     */
    $httpProvider.interceptors.unshift(function(){
      return {
       'request': function(config) {
          var addAuth = config.url.match(/_ah\/api\//) ||
                        config.url.match(/www\.googleapis\.com/);
          if (addAuth && authToken) {
            var auth = authToken.token_type + ' ' +
                       // TODO(alex): replace with id_token
                       authToken.access_token;
            config.headers['Authorization'] = auth;
          }
          return config;
        }
      };
    });

    this.$get = ApiFactory;
    /**
     * Our API client factory.
     */
    function ApiFactory($log, $q, $http, $rootScope, ApiRoot, ClientID, Scopes) {
      $rootScope.$on("GoogClientLoaded", function($event, gapiClient){
        gapi = gapiClient;
        login(true).always(function(){
          $rootScope.$broadcast('ApiReady');
        });
      });


      /**
       * Starts the auth flow in immediate or standar mode
       * @param  {boolean} immediate
       */
      function login(immediate) {
        var deferred = $q.defer(), promise = deferred.promise;

        promise.then(function(){
          $rootScope.$broadcast("Authenticated", true);
        }, function(err){
          $rootScope.$broadcast("Authenticated", false);
        })

        if (!gapi) {
          deferred.reject('Not ready yet');
          return promise;
        }

        if (authToken) {
          deferred.resolve(true);
        } else {
          gapi.auth.authorize({
            client_id: ClientID,
            scope: Scopes,
            immediate: immediate,
            response_type: 'token id_token'
          }, function(token){
            $log.debug('gapi.auth.authorize:', token);
            if (!token || token.error) {
              deferred.reject('auth with immediate=' + immediate + ' failed.');
            } else {
              authToken = token;
              gapi.auth.setToken(token);
              deferred.resolve(true);
            }
            $rootScope.$digest();
          });
        }

        return promise;
      }

      /**
       * Fetches a list of all messages submitted to the guestbook from the
       * backend.
       */
      function listGreetings(params) {
        var list = [];

        list.promise = $http.get(ApiRoot+'/greeting/v1/greetings', {
          params: params
        });

        list.promise.then(function(resp){
          list.length = 0;
          list.resolved = true;
          list.push.apply(list, resp.data.items);
        }, function(resp){
          $log.error(resp);
          list.error = resp.data.error_message ||
                       resp.data.error && resp.data.error.message ||
                       ('Error ' + resp.status);
        });

        return list;
      }

      /**
       * Submits a new message to the guestbook backend
       */
      function signGuestbook(greet) {
        return $http.post(ApiRoot+'/greeting/v1/greetings', greet);
      }

      /**
       * Fetches currently logged in user info (if any) from Google's OAuth2
       * endpoint.
       */
      function getUserInfo() {
        var user = {
          'promise': $http.get('https://www.googleapis.com/oauth2/v2/userinfo')
        };

        user.promise.then(
          function(resp){
            angular.copy(resp.data, user);
            user.resolved = true;
          }, function(resp){
            user.resolved = true;
            user.error = resp.data.error_message ||
                         resp.data.error && resp.data.error.message ||
                         ('Error ' + resp.status);
          });

        return user;
      }

      /**
       * Public-facing API object.
       */
      var api = {
        isLoggedIn: function() { return !!authToken },
        getToken: function() { return authToken },
        getUserInfo: getUserInfo,
        login: login,
        listGreetings: listGreetings,
        signGuestbook: signGuestbook
      }
      return api;
    }
  }


  // Global initialization

  /**
   * onGoogClientLoaded gets called when auth.js client is loaded
   * by the browser.
   */
  function onGoogClientLoaded() {
    var elem = window.document.querySelector("[ng-app]");
    angular.element(elem).scope().$emit("GoogClientLoaded", window.gapi);
  }

  /**
   * Exported functions.
   */
  window.GreetingCtrl = GreetingCtrl;
  window.onGoogClientLoaded = onGoogClientLoaded;
})(window);
