<?php
class SecretBase {
    private function hidden(): string { return 'private'; }
    protected static function guarded(): string { return 'protected'; }
    public function callbacks(): array {
        return [Closure::fromCallable([$this, 'hidden']), self::guarded(...)];
    }
}
$base = new SecretBase;
$callbacks = $base->callbacks();
$private = $callbacks[0];
$protected = $callbacks[1];
echo $private(), '|', $protected(), '|';
echo (int)is_callable([$base, 'hidden']), (int)is_callable(['SecretBase', 'guarded']);
