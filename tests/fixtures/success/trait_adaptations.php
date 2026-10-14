<?php
trait FirstName { public function name() { return "first"; } }
trait SecondName { public function name() { return "second"; } }
class Named {
    use FirstName, SecondName {
        FirstName::name insteadof SecondName;
        SecondName::name as private secondName;
    }
    public function both() { return $this->name() . "/" . $this->secondName(); }
}
echo (new Named)->both(), "\n";
