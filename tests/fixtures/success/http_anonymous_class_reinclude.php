<?php
class HttpAnonymousBase {
    public function value() { return 'base'; }
}
$first = require __DIR__ . '/../support/http_anonymous_class.inc';
$second = require __DIR__ . '/../support/http_anonymous_class.inc';
echo $first->value(), '|', $second->value(), '|', get_class($first) === get_class($second) ? 'same' : 'different';
