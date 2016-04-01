import DS from 'ember-data';

export default DS.Model.extend({
  type: DS.attr('string'),
  icon: function() {
    if (this.get('type') === 'directory') {
      return ('folder');
    }
    return ('description');
  }.property(),
  isSelected: false,
});
