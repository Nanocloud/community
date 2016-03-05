import OAuth2PasswordGrant from 'ember-simple-auth/authenticators/oauth2-password-grant';
import Ember from 'ember';

export default OAuth2PasswordGrant.extend({
  serverTokenEndpoint: 'oauth/token',

  clientId: '9405fb6b0e59d2997e3c777a22d8f0e617a9f5b36b6565c7579e5be6deb8f7ae:9050d67c2be0943f2c63507052ddedb3ae34a30e39bbbbdab241c93f8b5cf341',

  authenticate: function(username, password) {

    return new Ember.RSVP.Promise((resolve, reject) => {

      const serverTokenEndpoint = this.get('serverTokenEndpoint');
      const options = {
        url: serverTokenEndpoint,
        data: JSON.stringify({
          grant_type: 'password',
          username: username,
          password: password
        }),
        type: 'POST',
        dataType: 'json',
        contentType: 'application/json',
        headers: {
          Authorization: 'Basic ' + window.btoa(this.get('clientId'))
        }
      };

      Ember.$.ajax(options).then(function (response) {
          resolve(response);
      }, function (xhr) {
        reject(xhr.responseJSON || xhr.responseText);
      });
    });
  }
});
