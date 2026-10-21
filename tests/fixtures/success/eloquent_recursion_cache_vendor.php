<?php
require 'research/laravel-api/vendor/autoload.php';

class RecursionCacheProbe {
    use Illuminate\Database\Eloquent\Concerns\PreventsCircularRecursion;

    public int $calls = 0;

    public function run() {
        return $this->withoutRecursion(function () {
            ++$this->calls;
            if ($this->calls === 1) {
                return $this->run();
            }
            return 'later';
        }, 'guard');
    }

    public function fails() {
        return $this->withoutRecursion(function () {
            throw new Exception('stop');
        }, 'guard');
    }
}

$first = new RecursionCacheProbe;
$second = new RecursionCacheProbe;
echo $first->run(), '|', $first->run(), '|', $second->run(), '|';
echo $first->calls, ':', $second->calls, '|';
try {
    $first->fails();
} catch (Exception $e) {
    echo 'caught|';
}
try {
    $first->fails();
} catch (Exception $e) {
    echo 'caught';
}
