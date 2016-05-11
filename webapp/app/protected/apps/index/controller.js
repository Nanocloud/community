import Ember from 'ember';

export default Ember.Controller.extend({
  showSingleTab: false,
  showFileExplorer: false,
  connectionName: null,
  store: Ember.inject.service('store'),
  remoteSession: Ember.inject.service('remote-session'),
  session: Ember.inject.service('session'),
  isPublishing: false,

  applicationList: Ember.computed('model.@each', 'model.@each', function() {
    var array = this.get('model').rejectBy('alias', 'hapticPowershell');
    if (this.get('session.user.isAdmin') === true) {
      array = array.rejectBy('alias', 'hapticDesktop');
    }
    return array;
  }),

  actions: {

    disconnectGuacamole() {
      this.get('remoteSession').disconnectSession(this.get('connectionName'));
    },

    toggleSingleTab(connectionName) {
      this.set('connectionName', connectionName);
      this.toggleProperty('showSingleTab');
    },

    toggleFileExplorer() {
      this.toggleProperty('showFileExplorer');
    },
  }
});
