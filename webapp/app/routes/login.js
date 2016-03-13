import Ember from 'ember';

export default Ember.Route.extend({
  setupController: function(controller) {
    controller.set('identification', "");
    controller.set('password', "");
  }
});
