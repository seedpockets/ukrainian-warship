```
 _   _ _              _       _               _    _                _     _       
| | | | |            (_)     (_)             | |  | |              | |   (_)      
| | | | | ___ __ __ _ _ _ __  _  __ _ _ __   | |  | | __ _ _ __ ___| |__  _ _ __  
| | | | |/ / '__/ _` | | '_ \| |/ _` | '_ \  | |/\| |/ _` | '__/ __| '_ \| | '_ \ 
| |_| |   <| | | (_| | | | | | | (_| | | | | \  /\  / (_| | |  \__ \ | | | | |_) |
 \___/|_|\_\_|  \__,_|_|_| |_|_|\__,_|_| |_|  \/  \/ \__,_|_|  |___/_| |_|_| .__/ 
                                                                           | |    
                                                                           |_|  go birrr  
```
**Too fat to fight? Still want to do your part?**<br> 
Then we have just what you need! Strap your self in to you very own Ukrainian Warship and take aim!

Ukrainian Warship features:
- Acquire targets directly form IT ARMY of Ukraine
- Automatically detects online targets
- Stolen code from popular https stress test tool
- The pockets of dead target will gown sunflowers
<br>
<br>

# Install

##### Docker
```bash
docker pull sunflowerpockets/ukrainian-warship:latest
docker run -ti --rm sunflowerpockets/ukrainian-warship:latest
```

##### AWS EC2 User Data for Ubuntu Server 20.04 LTS
```bash
#!/bin/bash

apt update -y
apt upgrade -y
cd root
curl -OL https://go.dev/dl/go1.17.7.linux-amd64.tar.gz
tar -C /usr/local -xvf go1.17.7.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
git clone https://github.com/seedpockets/ukrainian-warship.git
cd ukrainian-warship/
mkdir bin
go build -ldflags "-s -w" -o /root/ukrainian-warship/bin/ukrainian-warship
chmod +x /root/ukrainian-warship/bin/ukrainian-warship
/root/ukrainian-warship/bin/ukrainian-warship kill > /dev/null
```

##### Ubuntu Server 20.04 LTS
```bash
#!/bin/bash

sudo apt update -y
sudo apt upgrade -y
curl -OL https://go.dev/dl/go1.17.7.linux-amd64.tar.gz
sudo tar -C /usr/local -xvf go1.17.7.linux-amd64.tar.gz
echo export PATH=$PATH:/usr/local/go/bin >> ~/.profile
source ~/.profile
git clone https://github.com/seedpockets/ukrainian-warship.git
cd ukrainian-warship/
mkdir bin
go build -ldflags "-s -w" -o bin/ukrainian-warship
chmod +x bin/ukrainian-warship
bin/ukrainian-warship kill    
```

Example output:
```
 Updates targets every 5 min...


Target
_________________________________________________________
https://www.nornickel.com/
https://rmk-group.ru/ru/
https://www.evraz.com/ru/
https://nangs.org/
https://www.metalloinvest.com/
https://www.polymetalinternational.com/ru/
https://www.sibur.ru/
https://www.uralkali.com/ru/
https://www.tmk-group.ru/
Total:  10
⣻
```