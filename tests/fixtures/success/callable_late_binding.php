<?php
class CallableParent {
    public static function who(): string { return static::class; }
    public static function saved(): Closure { return self::who(...); }
}
class CallableChild extends CallableParent {}
$callback = CallableChild::saved();
echo $callback(), '|', Closure::fromCallable('CallableChild::who')(), '|';
echo call_user_func(['CallableChild', 'who']);
