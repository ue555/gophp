<?php
require 'research/laravel-api/vendor/autoload.php';
class OnceObjectScope {
    public int $calls = 0;
    public function value() {
        return once(fn () => ++$this->calls);
    }
}
$first = new OnceObjectScope;
$second = new OnceObjectScope;
echo $first->value(), '|', $first->value(), '|', $second->value(), '|';
echo $first->calls, ':', $second->calls;
