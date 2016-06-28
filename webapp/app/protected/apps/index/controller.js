import Ember from 'ember';

export default Ember.Controller.extend({
  showSingleTab: false,
  showFileExplorer: false,
  connectionName: null,
  store: Ember.inject.service('store'),
  remoteSession: Ember.inject.service('remote-session'),
  session: Ember.inject.service('session'),
  configuration: Ember.inject.service('configuration'),
  isPublishing: false,

  applicationList: Ember.computed('model.@each', 'model.@each', function() {
    var array = this.get('model');
    if (this.get('session.user.isAdmin') === true) {
      array = array.rejectBy('alias', 'Desktop');
    }
    return array;
  }),

  launchVDI(connectionName) {
    return new Ember.RSVP.Promise((res) => {
      this.get('remoteSession').one('connected', () => {
        res();
      });

      this.set('connectionName', connectionName);
      this.set('showSingleTab', true);
    });
  },

  actions: {

    retryConnection() {
      this.toggleProperty('activator');
    },

    handleVdiClose() {
      this.get('remoteSession').disconnectSession(this.get('connectionName'));
    },

    toggleSingleTab(connectionName) {
      this.launchVDI(connectionName);
    },

    toggleFileExplorer() {
      this.toggleProperty('showFileExplorer');
    },
  }
});
