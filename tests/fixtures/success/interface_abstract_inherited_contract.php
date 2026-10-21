<?php
interface InheritedTask { public function execute($value, $extra = null); }
abstract class DeferredTask implements InheritedTask {}
class CompletedTask extends DeferredTask {
    public function execute($value, $extra = null) { return $value . ':' . $extra; }
}
echo (new CompletedTask)->execute('done', 'ok'), "\n";
