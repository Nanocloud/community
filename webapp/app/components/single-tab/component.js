import Ember from 'ember';

export default Ember.Component.extend({

  remoteSession: Ember.inject.service('remote-session'),

  connectionName: null,

  topBarItemToggleWindowCollector: {
    upload: false,
    clipboard: false,
    download: false,
  },

  showState: false,

  toggling() {

    if (this.get('showState') == false) {
      $('.canva-fullscreen').hide();
      $('.ember-modal-fullscreen').velocity({ opacity:1, left: 0} , {
        easing: "linear",
        duration: 400
      });

      setTimeout(function() {
          $('.canva-fullscreen').show();
      }.bind(this), 200);
      this.set('showState', true);
    }
    else {
      $('.canva-fullscreen').hide();
      $('.ember-modal-fullscreen').velocity({ opacity:0, left: -(window.innerWidth) }, {
        easing: "linear",
        duration: 400
      });

      setTimeout(function() {
        this.set('showState', false);
        this.set('isVisible', false);
        this.sendAction('onClose');
      }.bind(this), 900);
    }
  },

  initialize: function() {
    this.toggling();
  }.on('becameVisible'),

  uploadIsVisible: Ember.computed('topBarItemToggleWindowCollector.upload', function() {
    return this.get('topBarItemToggleWindowCollector.upload');
  }),

  clipboardIsVisible: Ember.computed('topBarItemToggleWindowCollector.clipboard', function() {
    return this.get('topBarItemToggleWindowCollector.clipboard');
  }),

  downloadIsVisible: Ember.computed('topBarItemToggleWindowCollector.download', function() {
    return this.get('topBarItemToggleWindowCollector.download');
  }),

  closeAll() {
    var object = this.get('topBarItemToggleWindowCollector');
    for (var prop in object) {
      var objToBeSet = 'topBarItemToggleWindowCollector.' + prop;
      this.set(objToBeSet, false);
    }
  },

  handleToggling(element) {
    var state = this.get('topBarItemToggleWindowCollector.' + element);
    this.closeAll();
    if (!state) {
      this.set('topBarItemToggleWindowCollector.' + element, true);
    }
  },

  actions: {

    toggleSingleTab() {
      this.toggling();
    },

    toggleUploadWindow() {
      this.handleToggling('upload');
    },

    toggleClipboardWindow() {
      this.handleToggling('clipboard');
    },

    toggleDownloadWindow() {
      this.handleToggling('download');
    },
  }
});
