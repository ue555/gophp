<?php
$date = DateTimeImmutable::createFromFormat('U.u', '1704164645.123456');
echo $date->format('Y-m-d\TH:i:s.u\Z'), "\n";
echo $date->format('Y|n|j|G|i|s|u|v|P|O|p|U.u'), "\n";
echo $date->setTimezone(new DateTimeZone('Asia/Tokyo'))->format('Y-m-d\TH:i:s.vP'), "\n";
echo $date->format('c'), "\n";
