import Ember from 'ember';

export default Ember.Route.extend({
  model() {
    console.log('refresh this2');
    return this.store.findAll('machine');
  }
});
