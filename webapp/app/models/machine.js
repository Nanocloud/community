import Ember from 'ember';
import DS from 'ember-data';

export default DS.Model.extend({
  name: DS.attr('string'),
  status: DS.attr('string'),
  ip: DS.attr('string'),
  adminPassword: DS.attr('string'),
  platform: DS.attr('string'),

  isUp: Ember.computed('status', function() {
    return this.get('status') === 'up';
  }),
  isDown: Ember.computed('status', function() {
    return this.get('status') === 'down';
  }),

  getPlatform: Ember.computed('platform', function() {
    switch (this.get('platform')) {
      case "vmwarefusion":
          return "VMware Fusion";
      case "qemu":
          return "Qemu";
      case "manual":
          return "Manual";
      default:
          return "No drivers detected";
    }
  }),
});
