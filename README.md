# Development paused:
### I'll come back to this project once I have a AD lab to test LDAP queries with and more time on my hands

# DropBox
Collection of scripts and software to help turn the Openwrt platform into a useful penetration testing tool

### Target Platform 
This will be intended for the 32-bit mips 24kc architecture.

### Status: 
- Most likely switching languages for the project to Rust for this 
- If I find a better low-power device, then I may switch the target architecture if it seems more ubiquitous. 
- Questioning effectiveness of using a matrix centered C2, might end up scraping that and going for HTTP2  or DNS if I am feeling up for that headache  
 

## Planned Tools

1. Dropbox: A remote access tool used for quickly issuing commands to the Openwrt Box, and logging data it gains from the internal network it was dropped on.
2. RogueBox: Standard automated rogue AP attacks for the Openwrt platform with phising capabilities.
3. EnumBox: A network enumeration tool which I still have not ironed out the details on.
