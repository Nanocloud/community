import Ember from 'ember';

export default Ember.Component.extend({

  remoteSession: Ember.inject.service('remote-session'),

  connectionName: null,

  topBarItemToggleWindowCollector: {
    upload: false,
    clipboard: false,
  },

  uploadIsVisible: Ember.computed('topBarItemToggleWindowCollector.upload', function() {
    return this.get('topBarItemToggleWindowCollector.upload');
  }),

  clipboardIsVisible: Ember.computed('topBarItemToggleWindowCollector.clipboard', function() {
    return this.get('topBarItemToggleWindowCollector.clipboard');
  }),

  closeAll() {
    var object = this.get('topBarItemToggleWindowCollector');
    for (var prop in object) {
      var objToBeSet = 'topBarItemToggleWindowCollector.' + prop;
      this.set(objToBeSet, false);
    }
  },

  handleToggling(element, state) {
    this.closeAll();
    if (!state) {
      this.set('topBarItemToggleWindowCollector.' + element, true);
    }
  },

  actions: {

    toggleSingleTab() {
      this.toggleProperty('isVisible');
    },

    toggleUploadWindow() {
      this.handleToggling('upload', this.get('topBarItemToggleWindowCollector.upload'));
    },

    toggleClipboardWindow() {
      this.handleToggling('clipboard', this.get('topBarItemToggleWindowCollector.clipboard'));
    },
  }
});
