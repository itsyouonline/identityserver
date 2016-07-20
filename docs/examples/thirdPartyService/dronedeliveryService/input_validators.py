
from wtforms.validators import ValidationError

def multiple_of(mult):
    ''' check if value is multipe of mult'''

    message = 'Must be multiple of %d' % (mult)

    def _multiple_of(form, field):
        if field.data % mult != 0:
            raise ValidationError(message)

    return _multiple_of
