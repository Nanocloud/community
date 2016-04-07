import Ember from 'ember';

export default Ember.Controller.extend({
  machines: Ember.computed('model.@each.isNew', 'model.@each.isDeleted', function() {
    return this.get('model').filterBy('isNew', false).filterBy('isDeleted', false);
  }),

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
