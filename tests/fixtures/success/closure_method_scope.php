<?php
class SecretBox {
    private $secret = 4;
    public function reader() { return function () { return $this->secret; }; }
}
$read = (new SecretBox)->reader();
echo $read(), "\n";
