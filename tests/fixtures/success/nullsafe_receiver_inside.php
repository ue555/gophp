<?php
class InnerReceiver { function take() { return "ok"; } }
class InsideBox {
    private $secret;
    protected $other;
    function __construct() { $this->secret = new InnerReceiver; $this->other = new InnerReceiver; }
    function run() { return $this?->secret?->take() . "|" . $this?->other?->take(); }
}
echo (new InsideBox)->run(), "\n";
