import Ember from 'ember';

export default Ember.Component.extend({

  session: Ember.inject.service('session'),
  flow: null,
  state: null,
  defaultState: " - ",

  initializeState:function() {
    this.set('state', this.get('defaultState'));
  }.on('init'),

  updateUploadStatus: function () {
    this.set('getUploadState', this.get('state') || this.get('defaultState'));
  }.observes('state'),

  showElement() {
    Ember.$('.drag-n-drop-area').css("opacity", "0.6");
  },

  hideElement() {
    Ember.$('.drag-n-drop-area').css("opacity", "0");
  },

  dragEnter() {
    this.showElement();
  },

  dragLeave() {
    this.hideElement();
  },

  drop() {
    this.hideElement();
  },

  didInsertElement() {

    this.flow = new window.Flow({
      target: '/upload',
      headers: { Authorization : "Bearer " + this.get('session.access_token') },
      singleFile: true
    });

    this.flow.assignDrop(this.element);

    this.flow.on('filesSubmitted', () => {
      this.flow.upload();
    });

    this.flow.on('complete', () => {
      if (!this.aborted) {
        this.set('state', "Completed");
      }
      else {
        this.set('state', "Aborted");
      }
      setTimeout(() => {
        this.set('state', null);
      }, 3000);
    });

    this.flow.on('uploadStart', () => {
      this.set('state', 0);
    });

    this.flow.on('error', () => {
        this.set('state', "Error");
    });

    this.flow.on('fileProgress', (flow) => {
      this.set('state', Math.floor(flow.progress() * 100));
      if (this.get('state') === 99) {
        this.set('state', 'Reassembling');
      }
    });
  },

  actions: {

    cancelDownload() {
      this.set('state', "Aborted");
      this.flow.cancel();
    }
  },

});
