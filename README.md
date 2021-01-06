# rastreamento_gorreios
Março, 2020<br>

### **Sobre o Bot**<br>
O Telebot (rastreamento_gorreios) utiliza uma API pública de rastreio (https://api.linketrack.com/track) para rastrear N encomendas de um usuário do telegram. Possuindo funcionalidades como:<br>
- Rastreamento automático de encomenda definido em um determinado intervalo;<br>
- Rastreio único de uma encomenda;<br>
- Remoção de encomendas configuradas para serem atualizadas automáticas.<br>

### **Sobre o sistema**<br>
<p> Este pequeno sistema foi criado por hobby (acompanhar minhas encomendas, claro) e para conhecer mais a utilização do serviço de Bot oferecido pelo Telegram (https://core.telegram.org/bots). <strong>HOBBY</strong><br></p>

### **Configurações**<br>
+ Criar um Bot pelo BotFather do telegram<br>
+ Configurar uma variável de ambiente com o nome: <strong>BOT_KEY_RASTGORREIOS</strong> com o valor da chave do Bot criado<br>
+ Configurar o tempo de rastreio automático (em minutos), nome da variável de ambiente: <strong>RASTGORREIOS_TRACK_INTERVAL</strong> [default: 5 minutos]<br>
+ Subir o MongoDB e configurar na variável de ambiente com o nome: <strong>MONGO_GORREIOS</strong> [default: mongodb://localhost:27017/]<br>

### **Comandos**<br>
+ /atualizar {codigo_rastreio} - Retorna o atual status da encomenda atualizado<br>
+ /acompanhe {codigo_rastreio} - Solicita para que o bot fique acompanhando o status da encomenda e notifique caso ela atualize<br>
+ /remova {codigo_rastreio} - Solicita que o bot remova um código do rastreio automático que já/não foi entregue<br>
+ /comandos - Lista todos os comandos disponíveis
