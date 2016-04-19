import Ember from 'ember';

export default Ember.Controller.extend({

  store: Ember.inject.service('store'),
  session: Ember.inject.service('session'),
  download: Ember.inject.service('download'),
  items: null,

  actions : {
    downloadFile: function(filename) {
      this.get('download').downloadFile(this.get('session.access_token'), filename);
    },
  }
});
