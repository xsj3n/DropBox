# DropBox
Collection of scripts and software to help turn the Openwrt platform into a useful penetration testing tool

### Target Platform 
This will be intended for the 32-bit mips 24kc architecture.

### Status: 
- Most likely switching languages for the project to Rust for this 
- If I find a better low-power device, then I may switch the target architecture if it seems more ubiquitous. 
- Questioning effectiveness of using a matrix centered C2, might end up scraping that and going for HTTP2  or DNS if I am feeling up for that headache  
- Finished with the initial build for my x64 PE packer & will begin serious work on DropBox in a few weeks. If anyone happens to see this, feel free to let me know of    some things you'd like to see implemented 

## Planned Tools

1. Dropbox: A remote access tool used for quickly issuing commands to the Openwrt Box, and logging data it gains from the internal network it was dropped on.
2. RogueBox: Standard automated rogue AP attacks for the Openwrt platform with phising capabilities.
3. EnumBox: A network enumeration tool which I still have not ironed out the details on.
