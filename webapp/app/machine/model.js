import Ember from 'ember';
import DS from 'ember-data';

export default DS.Model.extend({
  name: DS.attr('string'),
  status: DS.attr('string'),
  ip: DS.attr('string'),
  adminPassword: DS.attr('string'),
  platform: DS.attr('string'),
  progress: DS.attr('number'),
  machineSize: DS.attr('string'),
  type: DS.belongsTo('machine-type'),
  driver: DS.belongsTo('machine-driver'),

  isUp: Ember.computed('status', function() {
    return this.get('status') === 'up';
  }),
  isDown: Ember.computed('status', function() {
    return this.get('status') === 'down';
  }),
  isDownloading: Ember.computed('status', function() {
    return this.get('status') === 'creating';
  }),

  getPlatform: Ember.computed('platform', function() {
    switch (this.get('platform')) {
      case "vmwarefusion":
          return "VMware Fusion";
      case "qemu":
          return "Qemu";
      case "manual":
          return "Manual";
      case "aws":
          return "AWS";
      case "openstack":
          return "Openstack";
      default:
          return false;
    }
  }),

  driverDetected: Ember.computed('platform', function() {
    return this.get('getPlatform')? true : false;
  }),
});
