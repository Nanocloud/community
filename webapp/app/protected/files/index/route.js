import Ember from 'ember';

export default Ember.Route.extend({
  setupController(controller, model) {
    controller.set('items', model);
  },

  model() {
    return this.get('store').query('file', { filename: './' })
      .catch((err) => {
        if (err.errors.length === 1 && err.errors[0].code === "000007") {
          this.toast.warning("Cannot list files because Windows is not running");
          this.transitionTo('protected.files.nowindows');
        } else {
          return this.send("error", err);
        }
      });
  }
});
