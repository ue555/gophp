<?php
interface Limits { public const LOW = 2; }
class LimitBase implements Limits {
    protected const HIGH = 9;
    public function range() { return self::LOW . "/" . self::HIGH; }
}
class LimitChild extends LimitBase {}
echo Limits::LOW, '|', LimitChild::LOW, '|', (new LimitChild)->range(), "\n";
