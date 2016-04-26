import Model from 'ember-data/model';

export default Model.extend({

    Id: DS.attr('string'),
    sessionName: DS.attr('string'),
    username: DS.attr('string'),
    state: DS.attr('string'),
});
