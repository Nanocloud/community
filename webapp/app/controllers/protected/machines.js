import Ember from 'ember';

export default Ember.Controller.extend({
  _refreshLoop: null,
  activateRefreshLoop() {
    if (!this.get('_refreshLoop')) {
      let loop = window.setInterval(() => {
        let model = this.get('model');
        model.update();
      }, 3000);
      this.set('_refreshLoop', loop);
    }
  }
});
