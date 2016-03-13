import Ember from 'ember';

export default Ember.Component.extend({

  session: Ember.inject.service('session'),
  loadingFile: null,
  flow: null,
  aborted: false,

  showElement() {
    $('.element-active-state').css("opacity", "1");
  },

  hideElement() {
    $('.element-active-state').css("opacity", "0");
  },

  dragEnter(e) {
    this.showElement();
  },

  dragLeave(e) {
    this.hideElement();
  },

  drop() {
    this.hideElement();
  },

  actions: {

    cancelDownload() {
      this.set('aborted', true);
      this.flow.cancel();
    }
  },

  didInsertElement() {

    this.flow = new Flow({
      target: '/upload',
      headers: { Authorization : "Bearer " + this.get('session.data.authenticated.access_token') },
      singleFile: true
    });

    this.flow.assignDrop(this.element);

    this.flow.on('filesSubmitted', () => {
      this.flow.upload();
    });

    this.flow.on('complete', (event, flow) => {
      if (!this.aborted)
        this.set('onCompleteMessage', "Completed");
      else
        this.set('onCompleteMessage', "Aborted");
      this.set('loadingFile', false);
      setTimeout(() => {
        this.set('onCompleteMessage', null);
        this.set('aborted', false);
      }, 3000);
    });

    this.flow.on('uploadStart', (event, flow) => {
      this.set('loadingFile', 0);
    });

    this.flow.on('fileProgress', (flow, file) => {
      this.set('loadingFile', Math.floor(flow.progress() * 100));
    });
  }
});
