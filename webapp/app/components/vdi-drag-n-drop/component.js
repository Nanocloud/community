import Ember from 'ember';

export default Ember.Component.extend({

  session: Ember.inject.service('session'),
  flow: null,
  state: null,
  progress: null,
  show: false,

  showElement() {
    this.set('show', true);
  },

  hideElement() {
    this.set('show', false);
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
      if (this.get('state') !== 'Aborted') {
        this.toast.success("Your file has been uploaded successfully!");
        this.set('state', "Completed");
      }
      setTimeout(() => {
        this.set('state', null);
      }, 3000);
    });

    this.flow.on('error', () => {
        this.set('state', "Error");
    });

    this.flow.on('fileProgress', (flow) => {

      this.set('time', flow.timeRemaining());
      this.set('sizeUploaded', flow.sizeUploaded());
      this.set('progress', Math.floor(flow.progress() * 100));
      var progress = this.get('progress');
      if (progress >= 90 && progress <= 100) {
        this.set('state', 'Reassembling');
        this.set('progress', null);
      }
      else {
        this.set('state', null);
      }
    });
  },

  stopUpload() {
    this.toast.info("Your upload has been aborted successfully!");
    this.set('state', "Aborted");
    this.set('progress', null);
    this.flow.cancel();
  },

  actions: {
    cancelUpload() {
      this.stopUpload();
    }
  },
});
