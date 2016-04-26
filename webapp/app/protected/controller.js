import Ember from 'ember';
import config from 'nanocloud/config/environment';

export default Ember.Controller.extend({
  connectionName: null,

  name: config.APP.name,
  version: config.APP.version,

  isNotAdmin: Ember.computed('session.user', 'session.user', function() {
    if (this.get('session.user.isAdmin') === true) {
      return false;
    }
    return true; 
  }),

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
