import Ember from 'ember';
import ENV from 'frontend/config/environment';

/* global jQuery */

export default Ember.Controller.extend({

  token: null,
  guac_token: null,
  connectionName: null,

  getGuacToken: function() {
    jQuery.post(ENV.GUACAMOLE_URL + 'api/tokens', {
      access_token: this.get('token')
    }, function(post) {
      this.set('guac_token', post.authToken);
    }.bind(this));

  }.observes('token')
});
