<?php
class MagicBridge {
    public function __call($method, $arguments) {
        return $method . ':' . $arguments[0];
    }
    public static function __callStatic($method, $arguments) {
        return $method . ':' . $arguments[0];
    }
    public function __toString() { return 'bridge'; }
}
$bridge = new MagicBridge;
echo $bridge->missing(7), '|', MagicBridge::other(8), '|', $bridge;
