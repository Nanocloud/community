import Ember from 'ember';

export default Ember.Component.extend({
  classNames: ['machine'],
  machine: null,

  shouldEnableLightBulb: Ember.computed('machine.isUp', function() {
    if (this.get('machine.status') !== 'down') {
      return true;
    }
    return false;
  })

});
