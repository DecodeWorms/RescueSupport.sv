# RescueSupport.sv
The RescueSupport services provides aid for underprivileged people

# System Design architecture
Client(Mobile application, web browser) -> Api -> Backend Server(RescueSupport) -> Databases(MongoDB)

# Api(Backend Services)
- Onboarding OR Oauth
    - Registration(Company and User)
    - KYC
    - Login
    - Change Password
    - Update Password

- Manage Profile  
   - View Profile
   - Update Profile

- Aid Offer Management
   - Create aid
   - Update aid
   - Delete aid
   - View aid

- Communication and Coordination
  - Allow communications between supporter and recipient
  - Track the delivery process

- Resource Pledging
 - Allow supporters to pledge for resources(food , drink and clothes etc)
 - Track pledged resources and their fulfillment status

- Volunteering opportunities
  - List available volunteering opportunities

- Feedback and Support
  - Collect feedback from the supporter on their experience
  - Provide support channel for any issues or questions

- Notification and Alerts
  - Send notifications about new aid request, updates and events .
  - Alert supporters to urgent needs or emergency situations.

- Partnership and Collaborations
  - Facilitate partnership with Supporters and the other organizations
  - Manage Collaborative aid supports
