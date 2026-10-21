<?php
class HiddenCallable {
    private function run(): int { return 1; }
}
echo 'before|';
$callback = (new HiddenCallable)->run(...);
echo 'after';
