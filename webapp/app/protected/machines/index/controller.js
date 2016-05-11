import Ember from 'ember';

export default Ember.Controller.extend({

  machinesController: Ember.inject.controller('protected.machines'),
  drivers: Ember.computed.alias('machinesController.drivers'),

  machines: Ember.computed('model.@each.isNew', 'model.@each.isDeleted', function() {
    return this.get('model').filterBy('isNew', false).filterBy('isDeleted', false);
  }),

  driverName : Ember.computed(function() {
    return this.get('drivers').objectAt(0).id;
  }),

  isConfigurable : function() {
    return(this.get("driverName") !== "qemu" &&
           this.get("driverName") !== "manual" &&
           this.get("driverName") !== "vmwarefusion");
  }.property("driverName"),

  _refreshLoop: null,
  activateRefreshLoop() {
    if (!this.get('_refreshLoop')) {
      let loop = window.setInterval(() => {
        let model = this.get('model');
        model.update();
      }, 3000);
      this.set('_refreshLoop', loop);
    }
  },

  actions: {
    downloadWindows: function() {
      let machine = this.store.createRecord('machine', {
        name: "windows-custom-server-127.0.0.1-windows-server-std-2012R2-amd64",
        adminPassword: "Nanocloud123+",
      });

      machine.save()
      .then(() => {
        this.toast.info("Windows is downloading");
      })
      .catch(() => {
        this.toast.error("Could not download Windows, please try again");
      });
    }
  }
});
