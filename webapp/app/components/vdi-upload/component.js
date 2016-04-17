import VdiWindowComponent from 'nanocloud/components/vdi-window/component';

export default VdiWindowComponent.extend({

  actions: {
    clearList() {
      this.sendAction('flushHistory');
    },
  }
});

