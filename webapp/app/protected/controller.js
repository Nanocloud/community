import Ember from 'ember';
import config from 'nanocloud/config/environment';

export default Ember.Controller.extend({
  connectionName: null,

  session: Ember.inject.service('session'),
  name: config.APP.name,
  version: config.APP.version,

  showSidebar: false,
  actions: {
    toggleSidebar() {
      this.toggleProperty('showSidebar');
    },

    logout() {
      this.get('session').invalidate();
    }
  }
});
