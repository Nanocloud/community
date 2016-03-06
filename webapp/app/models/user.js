import DS from 'ember-data';

export default DS.Model.extend({
  email: DS.attr('string'),
  is_activated: DS.attr('boolean'),
  is_admin: DS.attr('boolean'),
  first_name: DS.attr('string'),
  last_name: DS.attr('string'),
  sam: DS.attr('string'),
  windows_password: DS.attr('string'),
  password: DS.attr('string'),
  password2: DS.attr('string')
});
